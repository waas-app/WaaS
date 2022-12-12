import { observable } from 'mobx';
import { InfoRes } from './sdk/server_pb.d';

class GlobalAppState {
  @observable
  info?: InfoRes.AsObject;
}

export const AppState = new GlobalAppState();

console.info('see global app state by typing "window.AppState"');

Object.assign(window as any, {
  get AppState() {
    console.log('AppState', AppState)
    return JSON.parse(JSON.stringify(AppState));
  },
});
