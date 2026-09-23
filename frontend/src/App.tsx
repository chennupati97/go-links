import { useEffect, useState } from "react";
import { getLinks } from "./api";
import { Link } from "./types";
import CreateLinkForm from "./components/CreateLinkForm";
import LinkTable from "./components/LinkTable";

function App() {
  const [links, setLinks] = useState<Link[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  async function loadLinks() {
    try {
      setLoading(true);

      const data = await getLinks();

      setLinks(data);
      setError("");
    } catch {
      setError("Failed to load links.");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadLinks();
  }, []);

  return (
    <main>
      <h1>🚀 Go Links</h1>
      
      <p>Create and manage internal URL shortcuts.</p>

      <CreateLinkForm onLinkCreated={loadLinks} />

      {loading && <p>Loading...</p>}

      {error && <p>{error}</p>}

      {!loading && !error && (
        <LinkTable links={links} />
      )}
    </main>
  );
}

export default App;