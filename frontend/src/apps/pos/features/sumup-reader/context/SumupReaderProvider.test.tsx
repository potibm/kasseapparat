import { beforeEach, describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import SumupReaderProvider from "./SumupReaderProvider";
import SumupReaderHint from "../components/SumupReaderHint";
import { useSumupReader } from "../hooks/useSumupReader";
import * as ConfigHookModule from "@core/config/hooks/useConfig";
import * as PosApiModule from "@pos/api/usePosApi";
import * as ReaderStorage from "@core/localstorage/helper/local-storage-reader";

vi.mock("@core/config/hooks/useConfig", () => ({ useConfig: vi.fn() }));
vi.mock("@pos/api/usePosApi", () => ({ usePosApi: vi.fn() }));
vi.mock("@core/localstorage/helper/local-storage-reader", () => ({
  getCurrentReaderId: vi.fn(),
  setCurrentReaderId: vi.fn(),
}));

const readers = [
  { id: "rdr_1", name: "Bar 1", status: "paired" },
  { id: "rdr_2", name: "Entrance", status: "expired" },
];

const fetchSumupReaders = vi.fn();
let storedReader: string | undefined;

const Selected = () => {
  const { selectedReader } = useSumupReader();
  return <span data-testid="selected">{selectedReader?.name ?? "-"}</span>;
};

const renderProvider = () =>
  render(
    <SumupReaderProvider>
      <SumupReaderHint />
      <Selected />
    </SumupReaderProvider>,
  );

describe("SumupReaderProvider", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    storedReader = undefined;
    vi.mocked(ReaderStorage.getCurrentReaderId).mockImplementation(
      () => storedReader,
    );
    vi.mocked(ReaderStorage.setCurrentReaderId).mockImplementation((id) => {
      storedReader = id;
    });
    fetchSumupReaders.mockResolvedValue(readers);
    vi.mocked(PosApiModule.usePosApi).mockReturnValue({
      fetchSumupReaders,
    } as unknown as ReturnType<typeof PosApiModule.usePosApi>);
    vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
      sumupEnabled: true,
    } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);
  });

  it("hints at a missing reader and lets the cashier pick one", async () => {
    renderProvider();

    const hint = await screen.findByTestId("sumup-reader-hint");
    expect(hint).toHaveTextContent("No card reader selected");

    fireEvent.click(hint);

    fireEvent.click(await screen.findByTestId("sumup-reader-rdr_1"));

    expect(ReaderStorage.setCurrentReaderId).toHaveBeenCalledWith("rdr_1");
    await waitFor(() =>
      expect(screen.getByTestId("selected")).toHaveTextContent("Bar 1"),
    );
    expect(screen.queryByTestId("sumup-reader-hint")).not.toBeInTheDocument();
    expect(
      screen.queryByText("Card reader for this till"),
    ).not.toBeInTheDocument();
  });

  it("names a reader that cannot take payments", async () => {
    storedReader = "rdr_2";
    renderProvider();

    expect(await screen.findByTestId("sumup-reader-hint")).toHaveTextContent(
      "Card reader „Entrance“ is expired",
    );
  });

  it("shows nothing while the reader works", async () => {
    storedReader = "rdr_1";
    renderProvider();

    await waitFor(() =>
      expect(screen.getByTestId("selected")).toHaveTextContent("Bar 1"),
    );
    expect(screen.queryByTestId("sumup-reader-hint")).not.toBeInTheDocument();
  });

  it("says so when the readers cannot be loaded", async () => {
    fetchSumupReaders.mockRejectedValue(new Error("offline"));
    renderProvider();

    fireEvent.click(await screen.findByTestId("sumup-reader-hint"));

    expect(
      await screen.findByText("The card readers could not be loaded."),
    ).toBeInTheDocument();
  });

  it("explains how to pair a reader when there is none", async () => {
    fetchSumupReaders.mockResolvedValue([]);
    renderProvider();

    fireEvent.click(await screen.findByTestId("sumup-reader-hint"));

    expect(
      await screen.findByText(/No card reader is paired yet/),
    ).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: "Refresh" }));
    await waitFor(() => expect(fetchSumupReaders).toHaveBeenCalledTimes(3));
  });

  it("does nothing without SumUp", () => {
    vi.mocked(ConfigHookModule.useConfig).mockReturnValue({
      sumupEnabled: false,
    } as unknown as ReturnType<typeof ConfigHookModule.useConfig>);
    renderProvider();

    expect(fetchSumupReaders).not.toHaveBeenCalled();
    expect(screen.queryByTestId("sumup-reader-hint")).not.toBeInTheDocument();
  });
});
