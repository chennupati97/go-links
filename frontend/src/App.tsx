import { useEffect, useState } from "react";
import { fetchShortcuts } from "./api";
import { Shortcut } from "./types";
import RegisterShortcutForm from "./components/RegisterShortcutForm";
import ShortcutList from "./components/ShortcutList";

function App() {
  const [items, setItems] = useState<Shortcut[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [loadError, setLoadError] = useState("");

  async function refreshShortcuts() {
    try {
      setIsLoading(true);
      const data = await fetchShortcuts();
      setItems(data);
      setLoadError("");
    } catch {
      setLoadError("Could not load shortcuts.");
    } finally {
      setIsLoading(false);
    }
  }

  useEffect(() => {
    refreshShortcuts();
  }, []);

  return (
    <main>
      <h1>JumpAlias</h1>
      <p>Map short aliases to the URLs your team uses every day.</p>

      <RegisterShortcutForm onRegistered={refreshShortcuts} />

      {isLoading && <p>Loading shortcuts...</p>}
      {loadError && <p>{loadError}</p>}
      {!isLoading && !loadError && <ShortcutList items={items} />}
    </main>
  );
}

export default App;
