import { Sidebar } from "./Sidebar";
import { StatusBar } from "./StatusBar";

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <div className="flex h-screen flex-col overflow-hidden bg-background text-foreground">
      <div className="flex flex-1 overflow-hidden">
        <Sidebar />
        <main className="flex-1 overflow-auto">{children}</main>
      </div>
      <StatusBar />
    </div>
  );
}
