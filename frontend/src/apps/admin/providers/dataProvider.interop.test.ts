import { describe, it, expect } from "vitest";
import { parse } from "query-string";
import packageJson from "../../../../package.json";

/**
 * `query-string@7` is CommonJS and calls `require('decode-uri-component')` as a
 * function. `decode-uri-component@0.5.0` is ESM-only (`"type": "module"`), so
 * requiring it yields the namespace `{ __esModule, default }` rather than a
 * callable. Vite's dep pre-bundler makes the same mistake, which crashed every
 * react-admin list view with "decodeComponent is not a function": ra-core calls
 * `parse()` from `query-string` to read list filters, sort and pagination from
 * the URL, so any sort, filter or page interaction took the list down.
 *
 * An npm `overrides` entry pinning `decode-uri-component` to `^0.5.0` forced
 * that ESM-only build onto the CJS consumer. This guards the dependency tree
 * rather than our own code, so an automated dependency bump cannot reintroduce
 * it silently.
 */
describe("decode-uri-component CJS interop", () => {
  it("is not pinned by an npm override", () => {
    expect(packageJson.overrides ?? {}).not.toHaveProperty(
      "decode-uri-component",
    );
  });

  it("lets query-string decode a query string", () => {
    expect(() => parse("sort=createdAt&order=DESC&page=1")).not.toThrow();
  });
});
