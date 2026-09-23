export interface Shortcut {
  id: number;
  alias: string;
  destination: string;
  registeredAt: string;
}

export interface RegisterShortcutRequest {
  alias: string;
  destination: string;
}
