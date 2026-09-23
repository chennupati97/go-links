import { API_BASE_URL } from "../api";
import { Link } from "../types";

interface LinkTableProps {
  links: Link[];
}

function LinkTable({ links }: LinkTableProps) {
  if (links.length === 0) {
    return <p>No shortcuts have been created yet.</p>;
  }

  return (
    <section>
      <h2>Existing Go Links</h2>

      <table>
        <thead>
          <tr>
            <th>Shortcut</th>
            <th>Destination URL</th>
            <th>Action</th>
          </tr>
        </thead>

        <tbody>
          {links.map((link) => (
            <tr key={link.id}>
              <td>{link.slug}</td>

              <td>
                <a
                  href={link.url}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {link.url}
                </a>
              </td>

              <td>
                <a
                  href={`${API_BASE_URL}/go/${link.slug}`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Open
                </a>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}

export default LinkTable;
