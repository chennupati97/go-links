import { FormEvent, useState } from "react";
import { registerShortcut } from "../api";

interface RegisterShortcutFormProps {
  onRegistered: () => void;
}

function RegisterShortcutForm({ onRegistered }: RegisterShortcutFormProps) {
  const [alias, setAlias] = useState("");
  const [destination, setDestination] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [errorMessage, setErrorMessage] = useState("");

  async function onSubmit(event: FormEvent) {
    event.preventDefault();
    setErrorMessage("");

    if (!alias.trim() || !destination.trim()) {
      setErrorMessage("Alias and destination URL are required.");
      return;
    }

    try {
      setSubmitting(true);
      await registerShortcut({
        alias: alias.trim(),
        destination: destination.trim(),
      });
      setAlias("");
      setDestination("");
      onRegistered();
    } catch (err: any) {
      if (err.response?.status === 409) {
        setErrorMessage("That alias is already taken.");
      } else if (err.response?.status === 400) {
        setErrorMessage("Please provide a valid http(s) destination.");
      } else {
        setErrorMessage("Could not register the shortcut.");
      }
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <section>
      <h2>Register a shortcut</h2>

      <form onSubmit={onSubmit}>
        <div>
          <label htmlFor="alias">Alias</label>
          <input
            id="alias"
            type="text"
            value={alias}
            onChange={(e) => setAlias(e.target.value)}
            placeholder="design-system"
          />
        </div>

        <div>
          <label htmlFor="destination">Destination URL</label>
          <input
            id="destination"
            type="url"
            value={destination}
            onChange={(e) => setDestination(e.target.value)}
            placeholder="https://example.com"
          />
        </div>

        {errorMessage && <p>{errorMessage}</p>}

        <button type="submit" disabled={submitting}>
          {submitting ? "Saving..." : "Save shortcut"}
        </button>
      </form>
    </section>
  );
}

export default RegisterShortcutForm;
