import React from 'react';
import { BrowserRouter,Routes, Route } from 'react-router-dom';

import Bingo from './components/Bingo';
import Qr from './components/Qr';
import Manual from './components/Manual';
import Home from './components/Home';
import Prize from './components/Prize';

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/bingo" element={<Bingo />} />
        <Route path="/qr" element={<Qr />} />
        <Route path="/manual" element={<Manual />} />
        <Route path="/prize" element={<Prize />} />
        <Route path="/" element={<Home />} />
      </Routes>
    </BrowserRouter>
  );
};
export default App;