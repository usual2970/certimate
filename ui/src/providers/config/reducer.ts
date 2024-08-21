import { Access } from "@/domain/access";
import { ConfigData } from ".";

type Action =
  | { type: "ADD_ACCESS"; payload: Access }
  | { type: "DELETE_ACCESS"; payload: string }
  | { type: "UPDATE_ACCESS"; payload: Access }
  | { type: "SET_ACCESSES"; payload: Access[] };

export const configReducer = (
  state: ConfigData,
  action: Action
): ConfigData => {
  switch (action.type) {
    case "SET_ACCESSES": {
      return {
        ...state,
        accesses: action.payload,
      };
    }
    case "ADD_ACCESS": {
      return {
        ...state,
        accesses: [action.payload, ...state.accesses],
      };
    }
    case "DELETE_ACCESS": {
      return {
        ...state,
        accesses: state.accesses.filter(
          (access) => access.id !== action.payload
        ),
      };
    }
    case "UPDATE_ACCESS": {
      return {
        ...state,
        accesses: state.accesses.map((access) =>
          access.id === action.payload.id ? action.payload : access
        ),
      };
    }
    default:
      return state;
  }
};
