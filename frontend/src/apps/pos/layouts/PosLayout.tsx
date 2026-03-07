export const PosLayout = ({ children, sidebar, topAlert, overlays }) => (
  <div className="App p-2 dark:bg-black min-h-screen">
    {topAlert}
    <div className="flex w-full">
      <main className="w-9/12 overflow-hidden">{children}</main>
      <aside className="fixed inset-y-0 right-0 w-3/12 bg-slate-200 dark:bg-gray-900 p-2 text-center">
        {sidebar}
      </aside>
    </div>
    <div id="layout-overlays">{overlays}</div>
  </div>
);

export default PosLayout;
