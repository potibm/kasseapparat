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
import { SectionList } from "./SectionList";
import { SectionCreate } from "./SectionCreate";
import { SectionInput } from "./SectionInput";
import { SectionDeleteButton } from "./SectionDeleteButton";
import { SECTION_NOT_EMPTY_MESSAGE } from "./section.utils";

const sections = [
  { id: 1, name: "Tickets", pos: 1 },
  { id: 2, name: "", pos: 2 },
];

const authProviderFor = (role: string): AuthProvider => ({
  login: () => Promise.resolve(),
  logout: () => Promise.resolve(),
  checkAuth: () => Promise.resolve(),
  checkError: () => Promise.resolve(),
  getPermissions: () => Promise.resolve(role),
});

const dataProviderWith = (overrides: Partial<DataProvider> = {}) =>
  testDataProvider({
    getList: () =>
      Promise.resolve({ data: sections, total: sections.length }) as never,
    getMany: () => Promise.resolve({ data: sections }) as never,
    ...overrides,
  });

const renderWithAdmin = (
  ui: React.ReactElement,
  {
    role = "admin",
    dataProvider = dataProviderWith(),
  }: { role?: string; dataProvider?: DataProvider } = {},
) =>
  render(
    <MemoryRouter>
      <AdminContext
        dataProvider={dataProvider}
        authProvider={authProviderFor(role)}
      >
        <ResourceContextProvider value="sections">
          {ui}
          <Notification />
        </ResourceContextProvider>
      </AdminContext>
    </MemoryRouter>,
  );

describe("SectionList", () => {
  it("lists named sections and unnamed ones as row breaks", async () => {
    renderWithAdmin(<SectionList />);

    expect(await screen.findByText("Tickets")).toBeInTheDocument();
    expect(screen.getByText("Row break #2")).toBeInTheDocument();
  });

  it("offers delete to admins", async () => {
    renderWithAdmin(<SectionList />);

    await screen.findByText("Tickets");
    await waitFor(() =>
      expect(screen.getAllByLabelText("Delete")).toHaveLength(2),
    );
  });

  it("hides delete from other users", async () => {
    renderWithAdmin(<SectionList />, { role: "user" });

    await screen.findByText("Tickets");
    expect(screen.queryByLabelText("Delete")).not.toBeInTheDocument();
  });
});

describe("SectionDeleteButton", () => {
  it("explains why a section with products cannot be deleted", async () => {
    const remove = vi.fn(() =>
      Promise.reject(new HttpError("Conflict", 409, { error: "Conflict" })),
    );

    renderWithAdmin(
      <RecordContextProvider value={sections[0]}>
        <SectionDeleteButton />
      </RecordContextProvider>,
      { dataProvider: dataProviderWith({ delete: remove as never }) },
    );

    fireEvent.click(await screen.findByLabelText("Delete"));
    fireEvent.click(await screen.findByText("Confirm"));

    expect(
      await screen.findByText(SECTION_NOT_EMPTY_MESSAGE),
    ).toBeInTheDocument();
    expect(remove).toHaveBeenCalled();
  });
});

describe("SectionCreate", () => {
  it("renders the name and position fields", async () => {
    renderWithAdmin(<SectionCreate />);

    expect(await screen.findByLabelText("Name")).toBeInTheDocument();
    expect(screen.getByLabelText(/Pos/)).toBeInTheDocument();
  });
});

describe("SectionInput", () => {
  it("preselects the first section", async () => {
    renderWithAdmin(
      <Form>
        <SectionInput preselectFirst />
      </Form>,
    );

    expect(await screen.findByText("Tickets")).toBeInTheDocument();
  });

  it("keeps an existing section", async () => {
    renderWithAdmin(
      <Form record={{ id: 7, sectionId: 2 }}>
        <SectionInput preselectFirst />
      </Form>,
    );

    expect(await screen.findByText("Row break #2")).toBeInTheDocument();
    expect(screen.queryByText("Tickets")).not.toBeInTheDocument();
  });
});
