import React from 'react';
import Card from '@mui/material/Card';
import CardHeader from '@mui/material/CardHeader';
import CardContent from '@mui/material/CardContent';
import Avatar from '@mui/material/Avatar';
import WifiIcon from '@mui/icons-material/Wifi';
import WifiOffIcon from '@mui/icons-material/WifiOff';
import MenuItem from '@mui/material/MenuItem';
import numeral from 'numeral';
import { lastSeen } from '../Util';
import { IconMenu } from './IconMenu';
import { PopoverDisplay } from './PopoverDisplay';
import { grpc } from '../Api';
import { observer } from 'mobx-react';

interface Props {
  device: import('../sdk/devices_pb.d').Device;
  onRemove: () => void;
}

@observer
export class DeviceListItem extends React.Component<Props> {
  removeDevice = async () => {
    try {
      const del = new (await import('../sdk/devices_pb.d')).DeleteDeviceReq()
      del.setName(this.props.device.getName())
      await grpc.devices.deleteDevice(del, null);
      this.props.onRemove();
    } catch {
      window.alert('api request failed');
    }
  };

  render() {
    const device = this.props.device;
    return (
      <Card>
        <CardHeader
          title={device.getName()}
          avatar={
            <Avatar style={{ backgroundColor: device.getConnected() ? '#76de8a' : '#bdbdbd' }}>
              {/* <DonutSmallIcon /> */}
              {device.getConnected() ? <WifiIcon /> : <WifiOffIcon />}
            </Avatar>
          }
          action={
            <IconMenu>
              <MenuItem style={{ color: 'red' }} onClick={this.removeDevice}>
                Delete
              </MenuItem>
            </IconMenu>
          }
        />
        <CardContent>
          <table cellPadding="5">
            <tbody>
              {device.getConnected() && (
                <>
                  <tr>
                    <td>Endpoint</td>
                    <td>{device.getEndpoint()}</td>
                  </tr>
                  <tr>
                    <td>Download</td>
                    <td>{numeral(device.getTransmitBytes()).format('0b')}</td>
                  </tr>
                  <tr>
                    <td>Upload</td>
                    <td>{numeral(device.getReceiveBytes()).format('0b')}</td>
                  </tr>
                </>
              )}
              {!device.getConnected() && (
                <tr>
                  <td>Disconnected</td>
                </tr>
              )}
              <tr>
                <td>Last Seen</td>
                <td>{lastSeen(device.getLastHandshakeTime().toObject())}</td>
              </tr>
              <tr>
                <td>Public key</td>
                <td>
                  <PopoverDisplay label="show">{device.getPublicKey()}</PopoverDisplay>
                </td>
              </tr>
            </tbody>
          </table>
        </CardContent>
      </Card>
    );
  }
}
