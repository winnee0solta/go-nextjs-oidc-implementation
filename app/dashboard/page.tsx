import { cookies } from "next/headers";

export default async function Dashboard() {
  // cookies() is synchronous in server components
  const cookieStore = await cookies();

  const cookieHeader = cookieStore
    .getAll()
    .map(({ name, value }) => `${name}=${value}`)
    .join("; ");

  const res = await fetch("http://localhost:8080/me", {
    headers: {
      Cookie: cookieHeader,
    },
    credentials: "include",
    cache: "no-store",
  });

  if (!res.ok) {
    return (
      <div>
        You are not logged in. <a href="/">Login</a>
      </div>
    );
  }

  const user = await res.json();
  console.log("🚀 ~ Dashboard ~ user:", user);

  return (
    <div className="min-h-screen bg-black flex flex-col items-center justify-center p-8 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="bg-black rounded-lg shadow-[0_4px_15px_rgba(255,255,255,0.1)] max-w-md w-full p-10 flex flex-col gap-8 text-white">
        <h1 className="text-3xl font-semibold">Dashboard</h1>

        <div className="space-y-2">
          <p className="text-lg">
            Welcome, <span className="font-medium">{user.Name}</span>
          </p>
          <p>
            Email: <span className="font-mono text-sm">{user.Email}</span>
          </p>
        </div>

        <a
          href="http://localhost:8080/oidc/logout"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-block w-full text-center rounded-full bg-white text-black hover:bg-gray-200 transition-colors font-semibold text-sm sm:text-base py-3 sm:py-4 select-none"
        >
          Logout
        </a>
      </main>
    </div>
  );
}
