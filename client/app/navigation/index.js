import React from 'react';
import { connect } from 'react-redux';
import { BackHandler } from "react-native";
// import { Notifications } from 'expo';
import { addNavigationHelpers, StackNavigator, NavigationActions } from 'react-navigation';

import LoginScreen from '../screens/LoginScreen';
import MainTabNavigator from './MainTabNavigator';
import { addListener } from '../store/middleware';
// import registerForPushNotificationsAsync from '../api/registerForPushNotificationsAsync';

export const AppNavigator = StackNavigator(
  {
    Main: {
      screen: MainTabNavigator,
    },
    Login: {
      screen: LoginScreen,
    },
  },
  {
    navigationOptions: () => ({
      headerTitleStyle: {
        fontWeight: 'normal',
      },
    }),
  }
);

class AppNavigation extends React.Component {
  componentDidMount() {
    BackHandler.addEventListener("hardwareBackPress", this.onBackPress);
    // this._notificationSubscription = this._registerForPushNotifications();
  }

  componentWillUnmount() {
    BackHandler.removeEventListener("hardwareBackPress", this.onBackPress);
    // this._notificationSubscription && this._notificationSubscription.remove();
  }

  onBackPress = () => {
    const { dispatch, navigationState } = this.props;
    if (navigationState.stateForLoggedIn.index <= 1) {
      BackHandler.exitApp();
      return;
    }
    dispatch(NavigationActions.back());
    return true;
  };


  render() {
    const { navigationState, dispatch, isAuthenticated } = this.props;
    const state = isAuthenticated
      ? navigationState.stateForLoggedIn
      : navigationState.stateForLoggedOut;

    return (
      <AppNavigator navigation={addNavigationHelpers({ dispatch, state, addListener })} />
    );
  }

  // _registerForPushNotifications() {
  //   // Send our push token over to our backend so we can receive notifications
  //   // You can comment the following line out if you want to stop receiving
  //   // a notification every time you open the app. Check out the source
  //   // for this function in api/registerForPushNotificationsAsync.js
  //   registerForPushNotificationsAsync();
  //
  //   // Watch for incoming notifications
  //   this._notificationSubscription = Notifications.addListener(this._handleNotification);
  // }

  // _handleNotification = ({ origin, data }) => {
  //   console.log(`Push notification ${origin} with data: ${JSON.stringify(data)}`);
  // };
}

const mapStateToProps = state => ({
  navigationState: state.navigation,
  isAuthenticated: state.userProfile.isAuthenticated
});

export default connect(
  mapStateToProps
)(AppNavigation);