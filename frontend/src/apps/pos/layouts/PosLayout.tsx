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
  topAlert?: ReactNode;
  overlays?: ReactNode;
}

export const PosLayout: React.FC<PosLayoutProps> = ({
  children,
  sidebar,
  topAlert,
  overlays,
}) => {
  useEffect(() => {
    const root = document.documentElement;
    root.classList.add(POS_TOUCH_CLASS);
    return () => root.classList.remove(POS_TOUCH_CLASS);
  }, []);

  return (
    <div className="App p-2 dark:bg-black min-h-screen">
      {topAlert}
      <div className="flex w-full">
        <main className="w-9/12 overflow-hidden">{children}</main>
        <aside className="fixed inset-y-0 right-0 w-3/12 min-w-50 bg-slate-200 dark:bg-gray-900 p-2 text-center">
          {sidebar}
        </aside>
      </div>
      <div id="layout-overlays">{overlays}</div>
    </div>
  );
};

export default PosLayout;
