import React from "react";
import { describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import {
  AdminContext,
  AuthProvider,
  DataProvider,
  Form,
  HttpError,
  Notification,
  RecordContextProvider,
  ResourceContextProvider,
  testDataProvider,
} from "react-admin";
import { MemoryRouter } from "react-router";
import { ProductGroupList } from "./ProductGroupList";
import { ProductGroupCreate } from "./ProductGroupCreate";
import { ProductGroupDeleteButton } from "./ProductGroupDeleteButton";
import { ProductLayoutInputs } from "./ProductLayoutInputs";
import { PRODUCT_GROUP_NOT_EMPTY_MESSAGE } from "./productGroup.utils";

const sections = [
  { id: 1, name: "Tickets", pos: 1 },
  { id: 2, name: "T-Shirts", pos: 2 },
];

const groups = [
  { id: 10, sectionId: 2, name: "Male", pos: 10 },
  { id: 11, sectionId: 2, name: "Female", pos: 20 },
];

const recordsOf = (resource: string) =>
  resource === "sections" ? sections : groups;

const authProviderFor = (role: string): AuthProvider => ({
  login: () => Promise.resolve(),
  logout: () => Promise.resolve(),
  checkAuth: () => Promise.resolve(),
  checkError: () => Promise.resolve(),
  getPermissions: () => Promise.resolve(role),
});

const dataProviderWith = (overrides: Partial<DataProvider> = {}) =>
  testDataProvider({
    getList: (resource: string) =>
      Promise.resolve({
        data: recordsOf(resource),
        total: recordsOf(resource).length,
      }) as never,
    getMany: (resource: string) =>
      Promise.resolve({ data: recordsOf(resource) }) as never,
    ...overrides,
  });

const renderWithAdmin = (
  ui: React.ReactElement,
  {
    role = "admin",
    resource = "productGroups",
    dataProvider = dataProviderWith(),
  }: { role?: string; resource?: string; dataProvider?: DataProvider } = {},
) =>
  render(
    <MemoryRouter>
      <AdminContext
        dataProvider={dataProvider}
        authProvider={authProviderFor(role)}
      >
        <ResourceContextProvider value={resource}>
          {ui}
          <Notification />
        </ResourceContextProvider>
      </AdminContext>
    </MemoryRouter>,
  );

describe("ProductGroupList", () => {
  it("lists the groups with their section", async () => {
    renderWithAdmin(<ProductGroupList />);

    // the first render of a list is slow on CI
    expect(
      await screen.findByText("Male", {}, { timeout: 5000 }),
    ).toBeInTheDocument();
    expect(screen.getByText("Female")).toBeInTheDocument();
    await waitFor(() =>
      expect(screen.getAllByText("T-Shirts").length).toBeGreaterThan(0),
    );
  });

  it("offers delete to admins only", async () => {
    const { unmount } = renderWithAdmin(<ProductGroupList />);

    await screen.findByText("Male");
    await waitFor(() =>
      expect(screen.getAllByLabelText(/delete/i)).toHaveLength(2),
    );
    unmount();

    renderWithAdmin(<ProductGroupList />, { role: "user" });

    await screen.findByText("Male");
    expect(screen.queryByLabelText(/delete/i)).not.toBeInTheDocument();
  });
});

describe("ProductGroupDeleteButton", () => {
  it("explains why a group with products cannot be deleted", async () => {
    const remove = vi.fn(() =>
      Promise.reject(new HttpError("Conflict", 409, { error: "Conflict" })),
    );

    renderWithAdmin(
      <RecordContextProvider value={groups[0]}>
        <ProductGroupDeleteButton />
      </RecordContextProvider>,
      { dataProvider: dataProviderWith({ delete: remove as never }) },
    );

    fireEvent.click(await screen.findByLabelText(/delete/i));
    fireEvent.click(await screen.findByText(/confirm/i));

    expect(
      await screen.findByText(PRODUCT_GROUP_NOT_EMPTY_MESSAGE),
    ).toBeInTheDocument();
    expect(remove).toHaveBeenCalled();
  });
});

describe("ProductGroupCreate", () => {
  it("renders name, section and position", async () => {
    renderWithAdmin(<ProductGroupCreate />);

    expect(await screen.findByLabelText(/name/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/pos/i)).toBeInTheDocument();
    expect(await screen.findByText("Tickets")).toBeInTheDocument();
  });
});

describe("ProductLayoutInputs", () => {
  it("asks for the section of a product without a group", async () => {
    renderWithAdmin(
      <Form>
        <ProductLayoutInputs preselectSection />
      </Form>,
      { resource: "products" },
    );

    expect(await screen.findByText("Tickets")).toBeInTheDocument();
    expect(screen.queryByLabelText(/variant/i)).not.toBeInTheDocument();
  });

  it("asks for the variant of a product in a group", async () => {
    renderWithAdmin(
      <Form record={{ id: 5, sectionId: 2, productGroupId: 10, variant: "S" }}>
        <ProductLayoutInputs />
      </Form>,
      { resource: "products" },
    );

    expect(await screen.findByText("Male")).toBeInTheDocument();
    expect(screen.getByLabelText(/variant/i)).toHaveValue("S");
    expect(screen.queryByText("T-Shirts")).not.toBeInTheDocument();
  });
});
