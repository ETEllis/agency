import Foundation

// MARK: - IPC Protocol

enum IPCMessageType: String, Codable {
    case broadcast
    case approval
    case bulletin
    case vote
    case handshake
    case ping
    case pong
}

struct IPCMessage: Codable {
    let type: IPCMessageType
    let payload: Data?

    enum CodingKeys: String, CodingKey {
        case type, payload
    }

    init(type: IPCMessageType, payload: Data? = nil) {
        self.type = type
        self.payload = payload
    }
}

struct IPCHandshake: Codable {
    let orgId: String
    let clientType: String
}

struct IPCVote: Codable {
    let proposalId: String
    let approved: Bool
}

// MARK: - Domain Models

struct BroadcastMessage: Identifiable, Equatable {
    let id = UUID()
    let actorID: String
    let message: String
    let createdAt: Date

    var initials: String {
        let last = actorID.split(separator: ".").last.map(String.init) ?? actorID
        let upper = last.uppercased()
        let chars = Array(upper)
        if chars.count >= 2 { return String(chars.prefix(2)) }
        if chars.count == 1 { return String(repeating: chars[0], count: 2) }
        return "AG"
    }

    var roleColor: RoleColor {
        RoleColor.forActor(actorID)
    }
}

struct ApprovalItem: Identifiable, Equatable {
    let id = UUID()
    let proposalID: String
    let actorID: String
    let actionType: String
    let target: String
    let createdAt: Date
}

struct BulletinEntry: Identifiable, Equatable {
    let id = UUID()
    let actorID: String
    let directive: String
    let output: String
    let score: Double
    let provider: String
    let modelID: String
    let createdAt: Date

    var scoreColor: ScoreColor {
        if score >= 0.75 { return .good }
        if score >= 0.5  { return .fair }
        return .poor
    }
}

// MARK: - Derived types

enum RoleColor: CaseIterable {
    case blue, purple, teal, green, orange, pink

    static func forActor(_ id: String) -> RoleColor {
        var hash = 0
        for c in id.unicodeScalars { hash = hash &* 31 &+ Int(c.value) }
        let cases = RoleColor.allCases
        return cases[abs(hash) % cases.count]
    }
}

enum ScoreColor {
    case good, fair, poor
}

// MARK: - JSON Payloads (mirrors Go IPC structs)

struct IPCBroadcastPayload: Decodable {
    let actorId: String
    let message: String
    let createdAt: Int64
}

struct IPCApprovalPayload: Decodable {
    let proposalId: String
    let actorId: String
    let actionType: String
    let target: String
    let createdAt: Int64
}

struct IPCBulletinPayload: Decodable {
    let actorId: String
    let directive: String
    let output: String
    let score: Double
    let provider: String
    let modelId: String
    let createdAt: Int64
}
