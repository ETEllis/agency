import SwiftUI

struct OfficeView: View {
    @EnvironmentObject var connection: AgencyConnection
    @State private var selectedTab: Tab = .broadcasts

    enum Tab: String, CaseIterable {
        case broadcasts = "Broadcasts"
        case bulletin   = "Bulletin"
        case approvals  = "Approvals"
    }

    var body: some View {
        NavigationSplitView {
            SidebarView(selectedTab: $selectedTab)
                .environmentObject(connection)
        } detail: {
            Group {
                switch selectedTab {
                case .broadcasts: BubbleListView()
                case .bulletin:   BulletinView()
                case .approvals:  ApprovalView()
                }
            }
            .environmentObject(connection)
        }
        .toolbar {
            ToolbarItem(placement: .status) {
                ConnectionBadge()
                    .environmentObject(connection)
            }
        }
    }
}

// MARK: - Sidebar

struct SidebarView: View {
    @EnvironmentObject var connection: AgencyConnection
    @Binding var selectedTab: OfficeView.Tab

    var body: some View {
        List(selection: $selectedTab) {
            Section("Office") {
                Label {
                    HStack {
                        Text("Broadcasts")
                        Spacer()
                        if !connection.broadcasts.isEmpty {
                            Text("\(connection.broadcasts.count)")
                                .font(.caption2)
                                .foregroundStyle(.secondary)
                        }
                    }
                } icon: {
                    Image(systemName: "bubble.left.and.bubble.right")
                }
                .tag(OfficeView.Tab.broadcasts)

                Label {
                    HStack {
                        Text("Bulletin")
                        Spacer()
                        if !connection.bulletins.isEmpty {
                            Text("\(connection.bulletins.count)")
                                .font(.caption2)
                                .foregroundStyle(.secondary)
                        }
                    }
                } icon: {
                    Image(systemName: "chart.bar.doc.horizontal")
                }
                .tag(OfficeView.Tab.bulletin)

                Label {
                    HStack {
                        Text("Approvals")
                        Spacer()
                        if !connection.approvals.isEmpty {
                            Text("\(connection.approvals.count)")
                                .foregroundStyle(.white)
                                .font(.caption2.bold())
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(.orange, in: Capsule())
                        }
                    }
                } icon: {
                    Image(systemName: "checkmark.shield")
                }
                .tag(OfficeView.Tab.approvals)
            }
        }
        .navigationSplitViewColumnWidth(200)
    }
}

// MARK: - Connection badge

struct ConnectionBadge: View {
    @EnvironmentObject var connection: AgencyConnection

    var body: some View {
        HStack(spacing: 5) {
            Circle()
                .fill(connection.isConnected ? Color.green : Color.red)
                .frame(width: 7, height: 7)
            Text(connection.isConnected ? "Connected" : "Offline")
                .font(.caption)
                .foregroundStyle(.secondary)
        }
    }
}
