// Use both for send messages to other client and to receive message fromother client
export interface Message {
  room: string;
  user: string;
  message: string;
  createdAt: string;
}

// return new Message
export function NewMessage(
  room: string,
  user: string,
  message: string
): Message {
  return {
    room: room,
    user: user,
    message: message,
    createdAt: Date.now().toString(),
  };
}
