import React from 'react';
import { observable } from 'mobx';
import { observer } from 'mobx-react';
import Popover from '@mui/material/Popover';
import IconButton from '@mui/material/IconButton';
import InfoIcon from '@mui/icons-material/Info';

interface Props {
  children: React.ReactNode;
}

@observer
export class Info extends React.Component<Props> {
  @observable
  anchor?: HTMLElement;

  render() {
    return (
      <>
        <IconButton onClick={(event) => (this.anchor = event.currentTarget)}>
          <InfoIcon />
        </IconButton>
        <Popover
          open={!!this.anchor}
          anchorEl={this.anchor}
          onClose={() => (this.anchor = undefined)}
          anchorOrigin={{
            vertical: 'bottom',
            horizontal: 'center',
          }}
          transformOrigin={{
            vertical: 'top',
            horizontal: 'center',
          }}
        >
          <div style={{ padding: 16 }}>{this.props.children}</div>
        </Popover>
      </>
    );
  }
}
