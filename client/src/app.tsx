import React from "react";
import { Route, Routes } from "react-router-dom";
import { AnalyticsPage } from "./pages/analytics-page";

function App() {
  return (
    <Routes>
      <Route path="/" element={<AnalyticsPage />} />
    </Routes>
  );
}

export default App;
