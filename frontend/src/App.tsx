import React from 'react';
import './App.css';
import { BrowserRouter } from 'react-router-dom';
import ResponsiveDrawer from './root/root';

function App() {
  return (
    <div className="App">
      <BrowserRouter>
        <ResponsiveDrawer />
      </BrowserRouter>
    </div>
  );
}

export default App;
