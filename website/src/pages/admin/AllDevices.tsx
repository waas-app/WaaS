import React from 'react';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import { observer } from 'mobx-react';
import { grpc } from '../../Api';
import { lastSeen, lazy } from '../../Util';
import { confirm } from '../../components/Present';
import { AppState } from '../../AppState';
import { StringValue } from 'google-protobuf/google/protobuf/wrappers_pb';

@observer
export class AllDevices extends React.Component {
  devices = lazy(async () => {
    const res = await grpc.devices.listAllDevices(new (await import('../../sdk/devices_pb.d')).ListAllDevicesReq(), null);
    return res.getItemsList();
  });

  deleteDevice = async (device: import('../../sdk/devices_pb.d').Device.AsObject) => {
    if (await confirm('Are you sure?')) {
      const del = new (await import('../../sdk/devices_pb.d')).DeleteDeviceReq()
      del.setName(device.name)
      del.setOwner(new StringValue().setValue(device.owner))
      await grpc.devices.deleteDevice(del, null);
      await this.devices.refresh();
    }
  };

  render() {
    if (!this.devices.current) {
      return <p>loading...</p>;
    }

    const rows = this.devices.current;

    // show the provider column
    // when there is more than 1 provider in use
    // i.e. not all devices are from the same auth provider.
    return (
      <div style={{ display: 'grid', gridGap: 25, gridAutoFlow: 'row'}}>
        <Typography variant="h5" component="h5">
          Devices
        </Typography>
        <TableContainer>
          <Table stickyHeader>
            <TableHead>
              <TableRow>
                <TableCell>Owner</TableCell>
                <TableCell>Device</TableCell>
                <TableCell>Connected</TableCell>
                <TableCell>Last Seen</TableCell>
                <TableCell>Actions</TableCell>
              </TableRow>
            </TableHead>
            <TableBody>
              {rows.map((row, i) => (
                <TableRow key={i}>
                  <TableCell component="th" scope="row">
                    {row.getOwnerName() || row.getOwnerEmail() || row.getOwner()}
                  </TableCell>
                  <TableCell>{row.getName()}</TableCell>
                  <TableCell>{row.getConnected() ? 'yes' : 'no'}</TableCell>
                  <TableCell>{lastSeen(row.getLastHandshakeTime().toObject())}</TableCell>
                  <TableCell>
                    <Button variant="outlined" color="secondary" onClick={() => this.deleteDevice(row.toObject())}>
                      Delete
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
        <Typography variant="h5" component="h5">
          Server Info
        </Typography>
        <code>
          <pre>
          {JSON.stringify(AppState.info, null, 2)}

          </pre>
        </code>
      </div>
    );
  }
}
