import React, { ReactNode } from "react";
import { Card } from "flowbite-react";

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
      <Card className="max-w-sm ">
        <h5 className="text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
          <img
            src="/android-chrome-192x192.png"
            alt="Kasseapparat"
            className="align-text-top h-7 inline"
          />{" "}
          Kasseapparat
        </h5>

        <div className="my-3">
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
