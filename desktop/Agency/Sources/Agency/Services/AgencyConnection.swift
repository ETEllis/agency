import Foundation
import Network

/// Manages the Unix socket connection to the Agency IPC server.
/// Publishes decoded events as async streams. Thread-safe.
@MainActor
final class AgencyConnection: ObservableObject {

    @Published var broadcasts: [BroadcastMessage] = []
    @Published var approvals: [ApprovalItem] = []
    @Published var bulletins: [BulletinEntry] = []
    @Published var isConnected: Bool = false
    @Published var connectionError: String? = nil

    private let orgID: String
    private let socketPath: String
    private var connection: NWConnection?
    private var readBuffer = Data()

    private static let decoder = JSONDecoder()

    init(orgID: String = "default", baseDir: String? = nil) {
        self.orgID = orgID
        let dir = baseDir ?? FileManager.default.currentDirectoryPath
        self.socketPath = "\(dir)/.agency/ipc-\(orgID).sock"
    }

    // MARK: - Connection lifecycle

    func connect() {
        connectionError = nil
        let endpoint = NWEndpoint.unix(path: socketPath)
        let params = NWParameters()
        params.allowLocalEndpointReuse = true
        let conn = NWConnection(to: endpoint, using: params)
        self.connection = conn

        conn.stateUpdateHandler = { [weak self] state in
            Task { @MainActor in
                self?.handleStateChange(state)
            }
        }
        conn.start(queue: .global(qos: .userInitiated))
    }

    func disconnect() {
        connection?.cancel()
        connection = nil
        isConnected = false
    }

    func sendVote(proposalID: String, approved: Bool) {
        guard isConnected else { return }
        let vote = IPCVote(proposalId: proposalID, approved: approved)
        guard let payload = try? JSONEncoder().encode(vote) else { return }
        let msg = IPCMessage(type: .vote, payload: payload)
        send(msg)
        // Remove from local approvals immediately for responsive UI.
        approvals.removeAll { $0.proposalID == proposalID }
    }

    // MARK: - Private

    private func handleStateChange(_ state: NWConnection.State) {
        switch state {
        case .ready:
            isConnected = true
            sendHandshake()
            receiveLoop()
        case .failed(let err):
            isConnected = false
            connectionError = err.localizedDescription
            scheduleReconnect()
        case .cancelled:
            isConnected = false
        default:
            break
        }
    }

    private func sendHandshake() {
        let hs = IPCHandshake(orgId: orgID, clientType: "desktop")
        guard let payload = try? JSONEncoder().encode(hs) else { return }
        let msg = IPCMessage(type: .handshake, payload: payload)
        send(msg)
    }

    private func send(_ msg: IPCMessage) {
        guard let conn = connection else { return }
        guard let data = try? JSONEncoder().encode(msg) else { return }
        var line = data
        line.append(0x0A) // newline delimiter
        conn.send(content: line, completion: .idempotent)
    }

    private func receiveLoop() {
        connection?.receive(minimumIncompleteLength: 1, maximumLength: 65536) { [weak self] data, _, isComplete, error in
            guard let self else { return }
            if let data {
                Task { @MainActor in
                    self.readBuffer.append(data)
                    self.drainBuffer()
                }
            }
            if let error {
                Task { @MainActor in
                    self.connectionError = error.localizedDescription
                }
                return
            }
            if !isComplete {
                self.receiveLoop()
            }
        }
    }

    private func drainBuffer() {
        while let newlineRange = readBuffer.range(of: Data([0x0A])) {
            let lineData = readBuffer[readBuffer.startIndex..<newlineRange.lowerBound]
            readBuffer.removeSubrange(readBuffer.startIndex...newlineRange.lowerBound)
            handleLine(lineData)
        }
    }

    private func handleLine(_ data: Data) {
        guard let msg = try? Self.decoder.decode(IPCMessage.self, from: data) else { return }
        guard let payload = msg.payload else { return }

        switch msg.type {
        case .broadcast:
            if let p = try? Self.decoder.decode(IPCBroadcastPayload.self, from: payload) {
                let bm = BroadcastMessage(
                    actorID: p.actorId,
                    message: p.message,
                    createdAt: Date(timeIntervalSince1970: Double(p.createdAt) / 1000)
                )
                broadcasts.append(bm)
                if broadcasts.count > 200 { broadcasts.removeFirst(50) }
            }

        case .approval:
            if let p = try? Self.decoder.decode(IPCApprovalPayload.self, from: payload) {
                let item = ApprovalItem(
                    proposalID: p.proposalId,
                    actorID: p.actorId,
                    actionType: p.actionType,
                    target: p.target,
                    createdAt: Date(timeIntervalSince1970: Double(p.createdAt) / 1000)
                )
                // Deduplicate.
                if !approvals.contains(where: { $0.proposalID == item.proposalID }) {
                    approvals.append(item)
                }
            }

        case .bulletin:
            if let p = try? Self.decoder.decode(IPCBulletinPayload.self, from: payload) {
                let entry = BulletinEntry(
                    actorID: p.actorId,
                    directive: p.directive,
                    output: p.output,
                    score: p.score,
                    provider: p.provider,
                    modelID: p.modelId,
                    createdAt: Date(timeIntervalSince1970: Double(p.createdAt) / 1000)
                )
                bulletins.append(entry)
                if bulletins.count > 500 { bulletins.removeFirst(100) }
            }

        case .pong:
            break // heartbeat acknowledged

        default:
            break
        }
    }

    private func scheduleReconnect() {
        DispatchQueue.main.asyncAfter(deadline: .now() + 3) { [weak self] in
            self?.connect()
        }
    }
}
