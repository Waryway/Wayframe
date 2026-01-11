import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ItemList } from './components/ItemList';
import { ItemDetail } from './components/ItemDetail';
import { CreateItem } from './components/CreateItem';
import './styles.css';

const App: React.FC = () => {
  return (
    <Router>
      <div className="app">
        <header className="app-header">
          <div className="container">
            <h1>Wayframe React Example</h1>
            <p>Full-stack Go + React application with REST API</p>
          </div>
        </header>

        <main className="app-main">
          <Routes>
            <Route path="/" element={<ItemList />} />
            <Route path="/items/:id" element={<ItemDetail />} />
            <Route path="/create" element={<CreateItem />} />
          </Routes>
        </main>

        <footer className="app-footer">
          <div className="container">
            <p>Powered by Wayframe</p>
          </div>
        </footer>
      </div>
    </Router>
  );
};

export default App;

