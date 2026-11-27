import React from "react";
import { Link } from "react-router";
import { HiCog, HiShoppingCart } from "react-icons/hi";
import BaseCard from "../../pos/components/BaseCard";

const tileClass =
  "flex min-h-16 items-center gap-3 rounded-lg border border-gray-200 px-4 text-lg font-semibold text-gray-900 hover:bg-gray-100 dark:border-gray-700 dark:text-white dark:hover:bg-gray-700";

/**
 * Start page with the two ways into Kasseapparat. The POS is made for the iPad;
 * on a phone the admin is the one that works, so both are offered here.
 */
export const Splash: React.FC = () => (
  <BaseCard title="">
    <p className="mb-4">A point of sale for events.</p>

    <nav aria-label="Start" className="flex flex-col gap-3">
      <Link to="/pos" className={tileClass}>
        <HiShoppingCart className="h-7 w-7 shrink-0" aria-hidden="true" />
        <span className="flex flex-col">
          POS
          <span className="text-sm font-normal text-gray-500 dark:text-gray-400">
            Selling at the till, made for the iPad
          </span>
        </span>
      </Link>
      {/* The admin has its own router, so it is loaded as a page of its own. */}
      <a href="/admin" className={tileClass}>
        <HiCog className="h-7 w-7 shrink-0" aria-hidden="true" />
        <span className="flex flex-col">
          Admin
          <span className="text-sm font-normal text-gray-500 dark:text-gray-400">
            Products, guest lists and purchases
          </span>
        </span>
      </a>
    </nav>
  </BaseCard>
);

export default Splash;
