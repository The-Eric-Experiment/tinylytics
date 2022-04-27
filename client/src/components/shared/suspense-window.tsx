import React, { FunctionComponent, Suspense } from "react";
import { Window, WindowContent, WindowHeader } from "react95";
import styled from "styled-components";
import { BaseWindow } from "./base-window";

interface SuspenseCardProps {
  title: string;
}

export const SuspenseWindow: FunctionComponent<SuspenseCardProps> = ({
  children,
  ...rest
}) => {
  return (
    <BaseWindow {...rest}>
      <Suspense fallback={<div>Loading...</div>}>{children}</Suspense>
    </BaseWindow>
  );
};
