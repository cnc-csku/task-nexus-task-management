'use client';

import { useSession } from "next-auth/react";

export default function Home() {
  const { data: session, status } = useSession();

  if (status === "loading") {
    return <div>Loading...</div>;
  }

  if (!session) {
    return <div>Access denied</div>;
  }

  return (
    <div>
      <pre>
        {JSON.stringify(session, null, 2)}
      </pre>
    </div>
  );
}
