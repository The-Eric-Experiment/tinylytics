import React, { Suspense } from "react";
import { Route, Routes } from "react-router-dom";
import { styleReset } from "react95";
// original Windows95 font (optionally)
// @ts-ignore
import ms_sans_serif from "react95/dist/fonts/ms_sans_serif.woff2";
// @ts-ignore
import ms_sans_serif_bold from "react95/dist/fonts/ms_sans_serif_bold.woff2";
// pick a theme of your choice
// @ts-ignore
import original from "react95/dist/themes/original";
import styled, { createGlobalStyle, ThemeProvider } from "styled-components";
import { AnalyticsPage } from "./pages/analytics-page";

const GlobalStyles = createGlobalStyle`
  @font-face {
    font-family: 'ms_sans_serif';
    src: url('${ms_sans_serif}') format('woff2');
    font-weight: 400;
    font-style: normal
  }
  @font-face {
    font-family: 'ms_sans_serif';
    src: url('${ms_sans_serif_bold}') format('woff2');
    font-weight: bold;
    font-style: normal
  }
  body {
    background-color: #018281;
    font-family: 'ms_sans_serif';
  }
  ${styleReset}
`;

function App() {
  return (
    <AppContainer>
      <GlobalStyles />
      <ThemeProvider theme={original}>
        <Suspense fallback={<div />}>
          <Routes>
            <Route path="/:domain" element={<AnalyticsPage />} />
            <Route path="/" element={<AnalyticsPage />} />
          </Routes>
        </Suspense>
      </ThemeProvider>
    </AppContainer>
  );
}

const AppContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

export default App;
