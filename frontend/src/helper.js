import chaturbate from './assets/images/logo_chaturbate.png';

export const eventTypes = [
  'tip', 'follow', 'unfollow', 'subscribe'
]

export const imageLookup = {
  chaturbate,
}

export const currencyName = {
  chaturbate: "tks",
}

export const genderLookup = {
  "m": "Male",
  "f": "Female",
  "t": "Trans",
  "c": "Couple",
}

export const subscriptionName = {
  chaturbate: {
    name: "fanclub",
    shortText: "joined fanclub",
    longText: "has joined your fanclub",
  },
}

export function getEventItemText(event) {
  switch (event.Type) {
    case "tip":
      const nativeTip = event.Event.TipValue + currencyName?.[event?.Provider]
      return {
        shortText: nativeTip,
        longText: `User ${event.Event.User.Username} has tipped ${nativeTip}`,
      }
    case "follow":
      return {
        shortText: "followed",
        longText: `User ${event.Event.User.Username} has followed you`,
      }
    case "unfollow":
      return {
        shortText: "unfollowed",
        longText: `User ${event.Event.User.Username} has unfollowed you`,
      }
    case "subscribe":
      return {
        shortText: "subscribed",
        longText: `User ${event.Event.User.Username} has ${subscriptionName[event.Provider].sentence}`,
      }
  }
}
export function capitalizeFirstLetter(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

