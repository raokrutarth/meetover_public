import React from 'react';
import { Platform } from 'react-native';
import { Ionicons } from '@expo/vector-icons';
import { TabNavigator, TabBarBottom } from 'react-navigation';

import Colors from '../constants/Colors';

import ListScreen from '../screens/ListScreen';
import MapScreen from '../screens/MapScreen';
import ChatsScreen from '../screens/ChatsScreen';

export default TabNavigator(
  {
    List: {
      screen: ListScreen,
    },
    Map: {
      screen: MapScreen,
    },
    Chats: {
      screen: ChatsScreen,
    },
  },
  {
    navigationOptions: ({ navigation }) => ({
      tabBarIcon: ({ focused, tintColor }) => {
        const { routeName } = navigation.state;
        let iconName;
        switch (routeName) {
          case 'List':
            iconName =
              Platform.OS === 'ios' ? `ios-list${focused ? '' : '-outline'}` : 'md-list';
            break;
          case 'Map':
            iconName = Platform.OS === 'ios' ? `ios-map${focused ? '' : '-outline'}` : 'md-map';
            break;
          case 'Chats':
            iconName =
              Platform.OS === 'ios' ? `ios-chatboxes${focused ? '' : '-outline'}` : 'md-chatboxes';
        }
        return (
          <Ionicons
            name={iconName}
            size={30}
            color={tintColor}
          />
        );
      },
    }),
    tabBarOptions: {
      activeTintColor: Colors.tabIconSelected,
    },
    tabBarComponent: TabBarBottom,
    tabBarPosition: 'bottom',
    animationEnabled: false,
    swipeEnabled: false,
  }
);
