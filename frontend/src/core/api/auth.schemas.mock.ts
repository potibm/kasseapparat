import { faker } from "@faker-js/faker";
import { LoginResponse, UserData } from "./auth.schemas";

/**
 * Generates a realistic mock UserData object for testing purposes.
 *
 * @param overrides  Specific fields to override in the generated object.
 */
export const createMockUserData = (overrides?: Partial<UserData>): UserData => {
  return {
    role: faker.helpers.arrayElement(["user", "admin"]),
    username: faker.internet.displayName(),
    gravatarUrl: faker.image.url(),
    id: faker.number.int({ min: 1, max: 99999 }),
    ...overrides,
  };
};

/**
 * Generates a realistic mock LoginResponse object for testing purposes.
 *
 * @param overrides  Specific fields to override in the generated object.
 */
export const createMockLoginResponse = (
  overrides?: Partial<LoginResponse>,
): LoginResponse => {
  return {
    access_token: faker.internet.jwt(),
    token_type: "Bearer",
    expires_in: faker.number.int({ min: 3600, max: 7200 }),
    role: faker.helpers.arrayElement(["user", "admin"]),
    username: faker.internet.displayName(),
    gravatarUrl: faker.image.url(),
    id: faker.number.int({ min: 1, max: 99999 }),
    ...overrides,
  };
};
