import { FormEvent, useState } from "react";
import { createLink } from "../api";

interface CreateLinkFormProps {
  onLinkCreated: () => void;
}

function CreateLinkForm({ onLinkCreated }: CreateLinkFormProps) {
  const [slug, setSlug] = useState("");
  const [url, setUrl] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();

    setError("");

    if (!slug.trim() || !url.trim()) {
      setError("Shortcut and destination URL are required.");
      return;
    }

    try {
      setLoading(true);

      await createLink({
        slug: slug.trim(),
        url: url.trim(),
      });

      setSlug("");
      setUrl("");

      onLinkCreated();
    } catch (err: any) {
      if (err.response?.status === 409) {
        setError("Shortcut already exists.");
      } else if (err.response?.status === 400) {
        setError("Please enter a valid URL.");
      } else {
        setError("Unable to create link.");
      }
    } finally {
      setLoading(false);
    }
  }

  return (
    <section>
      <h2>Create New Link</h2>

      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="slug">Shortcut</label>
          <input
            id="slug"
            type="text"
            value={slug}
            onChange={(e) => setSlug(e.target.value)}
            placeholder="design-system"
          />
        </div>

        <div>
          <label htmlFor="url">Destination URL</label>
          <input
            id="url"
            type="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="https://example.com"
          />
        </div>

        {error && <p>{error}</p>}

        <button type="submit" disabled={loading}>
          {loading ? "Creating..." : "Create Link"}
        </button>
      </form>
    </section>
  );
}

export default CreateLinkForm;
