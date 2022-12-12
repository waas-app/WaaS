import React from 'react';
import CssBaseline from '@mui/material/CssBaseline';
import Box from '@mui/material/Box';
import Navigation from './components/Navigation';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { observer } from 'mobx-react';
import { grpc } from './Api';
import { AppState } from './AppState';
import { YourDevices } from './pages/YourDevices';
import { AllDevices } from './pages/admin/AllDevices';
import Login  from './pages/auth/Login';

@observer
export class App extends React.Component {
  async componentDidMount() {
    AppState.info = await (await grpc.server.info(new (await import('./sdk/server_pb.d')).InfoReq(), null)).toObject();
  }

  render() {
    if (!AppState.info) {
      return <p>loading...</p>;
    }
    return (
      <Router>
        <CssBaseline />
        <Navigation />
        <Box component="div" m={2}>
          <Routes>
            <Route path="/" element={<YourDevices />} />
            {AppState.info.isAdmin && (
              <>
                <Route path="/admin/all-devices" element={<AllDevices />} />
              </>
            )}
            <Route path="/auth/login" element={<Login />} />
          </Routes>
        </Box>
      </Router>
    );
  }
}
