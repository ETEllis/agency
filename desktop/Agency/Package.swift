// swift-tools-version: 5.9
import PackageDescription

let package = Package(
    name: "Agency",
    platforms: [.macOS(.v14)],
    products: [
        .executable(name: "Agency", targets: ["Agency"]),
    ],
    targets: [
        .executableTarget(
            name: "Agency",
            path: "Sources/Agency",
            swiftSettings: [
                .enableExperimentalFeature("StrictConcurrency"),
            ]
        ),
    ]
)
