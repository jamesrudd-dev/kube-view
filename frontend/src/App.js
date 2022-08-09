import React from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';
import ClusterSelectionMenu from './components/ClusterSelectionMenu';

function App() {
  return (
    <div className="App">
      <h1>Kube View</h1>
      <ClusterSelectionMenu />
    </div>
  );
}

export default App;
