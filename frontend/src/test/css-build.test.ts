import { describe, it, expect } from "vitest";
import { readFileSync, readdirSync, existsSync } from "fs";
import { resolve, join } from "path";

/**
 * These tests verify the frontend build output is correct.
 * They catch issues like Tailwind CSS not being compiled (missing PostCSS plugin)
 * which would result in raw @apply/@tailwind directives in the output CSS.
 *
 * Run `pnpm run build` before running these tests for them to be meaningful.
 */
describe("CSS build output", () => {
  const distDir = resolve(__dirname, "../../dist");
  const cssDir = join(distDir, "static/css");

  it("dist directory exists", () => {
    expect(existsSync(distDir)).toBe(true);
  });

  it("index.html exists and references CSS files", () => {
    const indexPath = join(distDir, "index.html");
    expect(existsSync(indexPath)).toBe(true);

    const html = readFileSync(indexPath, "utf-8");
    expect(html).toContain('<link href="/static/css/');
    expect(html).toContain('rel="stylesheet"');
  });

  it("CSS files exist in dist/static/css/", () => {
    expect(existsSync(cssDir)).toBe(true);

    const files = readdirSync(cssDir).filter((f) => f.endsWith(".css"));
    expect(files.length).toBeGreaterThan(0);
  });

  it("CSS contains no uncompiled Tailwind directives", () => {
    const files = readdirSync(cssDir).filter((f) => f.endsWith(".css"));

    for (const file of files) {
      const css = readFileSync(join(cssDir, file), "utf-8");

      // These directives should be compiled away by @tailwindcss/postcss
      expect(css).not.toContain("@apply ");
      expect(css).not.toContain("@tailwind ");
      expect(css).not.toContain("@custom-variant ");
      expect(css).not.toMatch(/@theme\s*\{/);
    }
  });

  it("CSS contains compiled Tailwind utility classes", () => {
    const files = readdirSync(cssDir).filter((f) => f.endsWith(".css"));
    const allCss = files.map((f) => readFileSync(join(cssDir, f), "utf-8")).join("\n");

    // Core utility classes that must be present after compilation
    expect(allCss).toContain(".flex");
    expect(allCss).toContain(".items-center");
    expect(allCss).toContain(".hidden");
  });

  it("CSS contains theme variables (custom properties)", () => {
    const files = readdirSync(cssDir).filter((f) => f.endsWith(".css"));
    const allCss = files.map((f) => readFileSync(join(cssDir, f), "utf-8")).join("\n");

    // App theme variables should be defined
    expect(allCss).toContain("--background");
    expect(allCss).toContain("--foreground");
    expect(allCss).toContain("--primary");
    expect(allCss).toContain("--border");
  });

  it("CSS file is non-trivial size (> 10KB compiled)", () => {
    const files = readdirSync(cssDir).filter((f) => f.endsWith(".css"));
    const totalSize = files.reduce((sum, f) => {
      return sum + readFileSync(join(cssDir, f)).byteLength;
    }, 0);

    // Compiled Tailwind CSS should be at least 10KB
    // If it's tiny (< 5KB), it's likely uncompiled directives only
    expect(totalSize).toBeGreaterThan(10000);
  });

  it("index.html references JS files", () => {
    const html = readFileSync(join(distDir, "index.html"), "utf-8");
    expect(html).toContain('<script defer src="/static/js/');
  });

  it("all CSS links in index.html resolve to real files", () => {
    const html = readFileSync(join(distDir, "index.html"), "utf-8");
    const cssLinks = [...html.matchAll(/href="(\/static\/css\/[^"]+)"/g)];

    expect(cssLinks.length).toBeGreaterThan(0);

    for (const [, href] of cssLinks) {
      // Strip leading / to make it relative to dist
      const filePath = join(distDir, href);
      expect(existsSync(filePath)).toBe(true);
    }
  });

  it("all JS links in index.html resolve to real files", () => {
    const html = readFileSync(join(distDir, "index.html"), "utf-8");
    const jsLinks = [...html.matchAll(/src="(\/static\/js\/[^"]+)"/g)];

    expect(jsLinks.length).toBeGreaterThan(0);

    for (const [, src] of jsLinks) {
      const filePath = join(distDir, src);
      expect(existsSync(filePath)).toBe(true);
    }
  });
});
