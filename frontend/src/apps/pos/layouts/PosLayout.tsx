import React, { ReactNode, useEffect } from "react";

/**
 * Class set on <html> while the POS is mounted. It switches on the touch rules in
 * pos-style.css (no double-tap zoom, no text selection, tap feedback). It lives on
 * <html> rather than on the layout so it also covers modals that portal to <body>.
 */
export const POS_TOUCH_CLASS = "pos-touch";

interface PosLayoutProps {
  children: ReactNode;
  sidebar: ReactNode;
  /** Bar above the main content (product grid). */
  header?: ReactNode;
  topAlert?: ReactNode;
  overlays?: ReactNode;
}

export const PosLayout: React.FC<PosLayoutProps> = ({
  children,
  sidebar,
  header,
  topAlert,
  overlays,
}) => {
  useEffect(() => {
    const root = document.documentElement;
    root.classList.add(POS_TOUCH_CLASS);
    return () => root.classList.remove(POS_TOUCH_CLASS);
  }, []);

  return (
    // The screen never scrolls as a whole: the header stays in view and only the
    // products scroll, like the cart in the sidebar.
    <div className="App flex h-dvh flex-col overflow-hidden p-2 dark:bg-black">
      {topAlert}
      <div className="flex min-h-0 w-full flex-1">
        <main className="mr-[360px] flex min-h-0 min-w-0 flex-1 flex-col">
          <div className="shrink-0">{header}</div>
          <div
            data-testid="pos-products-scroll"
            className="min-h-0 flex-1 overflow-y-auto overscroll-contain"
          >
            {children}
          </div>
        </main>
        <aside className="fixed inset-y-0 right-0 flex w-[360px] flex-col overflow-hidden bg-slate-200 p-2 text-center dark:bg-gray-900">
          {sidebar}
        </aside>
      </div>
      <div id="layout-overlays">{overlays}</div>
    </div>
  );
};

export default PosLayout;
