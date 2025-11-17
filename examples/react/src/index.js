import React from 'react';
import { createRoot } from 'react-dom/client';
import './index.css';

function App() {
  return <h1>Hello from React (Bazel build)!</h1>;
}

const root = createRoot(document.getElementById('root'));
root.render(<App />);

