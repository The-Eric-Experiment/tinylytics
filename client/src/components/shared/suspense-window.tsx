import React, { FunctionComponent, PropsWithChildren, Suspense } from "react";
import { BaseWindow } from "./base-window";

interface SuspenseCardProps {
  title: string;
}

export const SuspenseWindow: FunctionComponent<
  PropsWithChildren<SuspenseCardProps>
> = ({ children, ...rest }) => {
  return (
    <BaseWindow {...rest}>
      <Suspense fallback={<div>Loading...</div>}>{children}</Suspense>
    </BaseWindow>
  );
};
