import React, { FunctionComponent, Suspense } from "react";
import { Card } from "../shared/page-layout";

interface SuspenseCardProps {}

export const SuspenseCard: FunctionComponent<SuspenseCardProps> = ({
  children,
}) => {
  return (
    <Card>
      <Suspense fallback={<div>Loading...</div>}>{children}</Suspense>
    </Card>
  );
};
