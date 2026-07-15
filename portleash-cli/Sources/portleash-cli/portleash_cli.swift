import Foundation

#if canImport(FoundationNetworking)
import FoundationNetworking
#endif

@main
struct PortLeashCLI {
    static func main() {
        print("Connecting to PortLeash Daemon...")

        let url = URL(string: "http://127.0.0.1:4848/status")!
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