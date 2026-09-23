import { BACKEND_ORIGIN } from "../api";
import { Shortcut } from "../types";

interface ShortcutListProps {
  items: Shortcut[];
}

function ShortcutList({ items }: ShortcutListProps) {
  if (items.length === 0) {
    return <p>No shortcuts registered yet.</p>;
  }

  return (
    <section>
      <h2>Your shortcuts</h2>

      <table>
        <thead>
          <tr>
            <th>Alias</th>
            <th>Destination</th>
            <th>Open</th>
          </tr>
        </thead>

        <tbody>
          {items.map((item) => (
            <tr key={item.id}>
              <td>{item.alias}</td>
              <td>
                <a
                  href={item.destination}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  {item.destination}
                </a>
              </td>
              <td>
                <a
                  href={`${BACKEND_ORIGIN}/j/${item.alias}`}
                  target="_blank"
                  rel="noopener noreferrer"
                >
                  Jump
                </a>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}

export default ShortcutList;
