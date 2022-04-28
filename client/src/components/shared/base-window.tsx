import React, { FunctionComponent, PropsWithChildren } from "react";
import { Window, WindowContent, WindowHeader } from "react95";
import styled from "styled-components";

interface BaseWindowProps {
  title: string;
}

export const BaseWindow: FunctionComponent<
  PropsWithChildren<BaseWindowProps>
> = ({ title, children }) => {
  return (
    <WindowStyled className="window">
      <StyledWindowHeader className="window-header">
        <span>{title}</span>
      </StyledWindowHeader>
      <WindowContentStyled>{children}</WindowContentStyled>
    </WindowStyled>
  );
};

const StyledWindowHeader = styled(WindowHeader)`
  display: inline-flex;
  flex-direction: column;
  > span {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
`;

const WindowStyled = styled(Window)`
  height: 100%;
  display: flex;
  flex-direction: column;
`;

const WindowContentStyled = styled(WindowContent)`
  flex: 1;
`;
