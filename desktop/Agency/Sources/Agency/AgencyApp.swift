import SwiftUI

@main
struct AgencyApp: App {
    @StateObject private var connection: AgencyConnection = {
        let orgID = ProcessInfo.processInfo.environment["AGENCY_ORG_ID"] ?? "default"
        let baseDir = ProcessInfo.processInfo.environment["AGENCY_BASE_DIR"]
        return AgencyConnection(orgID: orgID, baseDir: baseDir)
    }()

    var body: some Scene {
        WindowGroup {
            OfficeView()
                .environmentObject(connection)
                .onAppear { connection.connect() }
                .onDisappear { connection.disconnect() }
        }
        .windowStyle(.hiddenTitleBar)
        .defaultSize(width: 1100, height: 740)
        .commands {
            CommandGroup(after: .appInfo) {
                Button("Reconnect") { connection.connect() }
                    .keyboardShortcut("r", modifiers: [.command, .shift])
            }
        }

        // Approval panel — secondary window, auto-pops when there are pending items.
        Window("Approvals", id: "approvals") {
            ApprovalView()
                .environmentObject(connection)
        }
        .defaultSize(width: 440, height: 480)
    }
}
