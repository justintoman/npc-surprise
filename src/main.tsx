import { Provider as JotaiProvider } from 'jotai';
import React from 'react';
import ReactDOM from 'react-dom/client';
import { store } from '~/state.ts';
import { ThemeProvider } from '~/ThemeProvider.tsx';
import { App } from './App.tsx';
import './index.css';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <JotaiProvider store={store}>
      <ThemeProvider>
        <App />
      </ThemeProvider>
    </JotaiProvider>
  </React.StrictMode>,
);
