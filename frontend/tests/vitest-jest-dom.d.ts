/* eslint-disable @typescript-eslint/no-empty-object-type */
/* eslint-disable @typescript-eslint/no-explicit-any */
import "vitest";
import { type TestingLibraryMatchers } from "@testing-library/jest-dom/types/matchers";

declare module "vitest" {
  interface Assertion<
    R extends void | Promise<void>,
    T,
  > extends TestingLibraryMatchers<R, T> {}

  interface AsymmetricMatchersContaining extends TestingLibraryMatchers<
    any,
    any
  > {}
}
