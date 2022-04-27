import React, { Suspense } from "react";
import { Route, Routes } from "react-router-dom";
import { styleReset } from "react95";
// original Windows95 font (optionally)
// pick a theme of your choice
// @ts-ignore
import original from "react95/dist/themes/original";
import styled, { createGlobalStyle, ThemeProvider } from "styled-components";
import { AnalyticsPage } from "./pages/analytics-page";

const GlobalStyles = createGlobalStyle`
  ${styleReset};
`;

function App() {
  return (
    <>
      <GlobalStyles />
      <AppContainer>
        <ThemeProvider theme={original}>
          <Suspense fallback={<div />}>
            <Routes>
              <Route path="/:domain" element={<AnalyticsPage />} />
              <Route path="/" element={<AnalyticsPage />} />
            </Routes>
          </Suspense>
        </ThemeProvider>
      </AppContainer>
    </>
  );
}

const AppContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

export default App;
