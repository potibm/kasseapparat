import React, { ReactNode } from "react";
import { Card } from "flowbite-react";
import Logo from "@core/logo/Logo";

interface MinimalBaseCardProps {
  children: ReactNode;
  title?: string;
  navigation?: ReactNode;
}

const MinimalBaseCard: React.FC<MinimalBaseCardProps> = ({
  children,
  title = null,
  navigation,
}) => {
  return (
    <div className="flex justify-center items-center h-screen">
      <Card className="max-w-sm w-full">
        <h5 className="flex items-center gap-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          <Logo className="h-7 w-7 text-[#ff3873]" />
          Kasseapparat
        </h5>

        <div className="my-3 dark:text-gray-200">
          {title && <h2 className="text-xl mb-2">{title}</h2>}

          {children}
        </div>
        {navigation && (
          <>
            <hr />
            <p className="text-xs dark:text-gray-200">{navigation}</p>
          </>
        )}
      </Card>
    </div>
  );
};

export default MinimalBaseCard;
