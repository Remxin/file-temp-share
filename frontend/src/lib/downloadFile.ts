import { getStoredToken } from "@/lib/storage/token";

const apiBase = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

type DownloadOptions = {
  path: string;
  filename: string;
  transferPassword?: string;
  useAuth?: boolean;
};

export async function downloadFile({
  path,
  filename,
  transferPassword,
  useAuth = false,
}: DownloadOptions): Promise<void> {
  const headers = new Headers();
  if (useAuth) {
    const token = getStoredToken();
    if (token) headers.set("Authorization", `Bearer ${token}`);
  }
  if (transferPassword) {
    headers.set("X-Transfer-Password", transferPassword);
  }

  const res = await fetch(`${apiBase}${path}`, { headers });
  if (!res.ok) {
    let message = "Download failed";
    try {
      const body = (await res.json()) as { error?: string };
      if (body.error) message = body.error;
    } catch {
      /* ignore */
    }
    throw new Error(message);
  }

  const blob = await res.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  a.remove();
  URL.revokeObjectURL(url);
}
