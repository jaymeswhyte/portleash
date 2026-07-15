import Foundation
import ArgumentParser

#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

@main
struct PortLeashCLI: ParsableCommand {
    @Argument(help: "Command (find)")
    var command: String

    @Argument(help: "Port to find")
    var port: Int? = nil

    func run() throws {
        print("Connecting to PortLeash Daemon...")

        let url = URL(string: "http://127.0.0.1:4848/status")!
        if command == "find"
        {
            print("Finding task for port \(port ?? 0)")
        }
        let task = URLSession.shared.dataTask(with: url) { data, response, error in
            if let error = error {
                print("Error connecting to daemon: \(error.localizedDescription)")
                return
            }
            
            guard let data = data else { return }
            if let jsonString = String(data: data, encoding: .utf8) {
                print("Daemon Response: \(jsonString)")
            }
        }
        task.resume()

        RunLoop.main.run(until: Date(timeIntervalSinceNow: 2))
    }
}