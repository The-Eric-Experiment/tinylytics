import React, { FunctionComponent, Suspense } from "react";
import { Window, WindowContent, WindowHeader } from "react95";
import styled from "styled-components";

interface SuspenseCardProps {
  title: string;
}

export const SuspenseWindow: FunctionComponent<SuspenseCardProps> = ({
  title,
  children,
}) => {
  return (
    <WindowStyled className="window">
      <WindowHeader className="window-header">
        <span>{title}</span>
      </WindowHeader>
      <WindowContentStyled>
        <Suspense fallback={<div>Loading...</div>}>{children}</Suspense>
      </WindowContentStyled>
    </WindowStyled>
  );
};

const WindowStyled = styled(Window)`
  height: 100%;
  display: flex;
  flex-direction: column;
`;

const WindowContentStyled = styled(WindowContent)`
  flex: 1;
`;
