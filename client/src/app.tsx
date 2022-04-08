import React, { Suspense } from "react";
import { Routes, Route } from "react-router-dom";
import { AnalyticsPage } from "./pages/analytics-page";

function App() {
  return (
    <Suspense fallback={<div />}>
      <Routes>
        <Route path="/:domain" element={<AnalyticsPage />} />
        <Route path="/" element={<AnalyticsPage />} />
      </Routes>
    </Suspense>
  );
}

export default App;
