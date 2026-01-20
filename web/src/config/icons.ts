import {
  DashboardOutlined,
  UserOutlined,
  SettingOutlined,
  TeamOutlined,
  CloudOutlined,
  FileOutlined,
  UploadOutlined,
  DatabaseOutlined,
  AppstoreOutlined,
  BookOutlined,
  BarsOutlined,
  SafetyOutlined,
  ToolOutlined,
  MonitorOutlined,
  EditOutlined,
  FormOutlined,
  FolderOutlined,
  FileTextOutlined,
  ApartmentOutlined,
  UnorderedListOutlined,
  LoginOutlined,
  FileSearchOutlined,
  CloudServerOutlined,
  FileImageOutlined,
  HomeOutlined,
  UserAddOutlined,
  UsergroupAddOutlined,
  SolutionOutlined,
  ProfileOutlined,
  ControlOutlined,
  FundOutlined,
  BarChartOutlined,
  LineChartOutlined,
  PieChartOutlined,
  ShoppingCartOutlined,
  ShopOutlined,
  TagOutlined,
  TagsOutlined,
  BellOutlined,
  MessageOutlined,
  MailOutlined,
  PhoneOutlined,
  CalendarOutlined,
  ClockCircleOutlined,
  EnvironmentOutlined,
  GlobalOutlined,
  SearchOutlined,
  FilterOutlined,
  SortAscendingOutlined,
  SortDescendingOutlined,
  ReloadOutlined,
  SyncOutlined,
  DownloadOutlined,
  ExportOutlined,
  ImportOutlined,
  PrinterOutlined,
  CopyOutlined,
  ScissorOutlined,
  DeleteOutlined,
  SaveOutlined,
  CheckOutlined,
  CloseOutlined,
  PlusOutlined,
  MinusOutlined,
  EyeOutlined,
  EyeInvisibleOutlined,
  LockOutlined,
  UnlockOutlined,
  KeyOutlined,
  ApiOutlined,
  CodeOutlined,
  BugOutlined,
  RocketOutlined,
  ThunderboltOutlined,
  FireOutlined,
  StarOutlined,
  HeartOutlined,
  LikeOutlined,
  DislikeOutlined,
  SmileOutlined,
  FrownOutlined,
  MehOutlined,
  QuestionCircleOutlined,
  InfoCircleOutlined,
  WarningOutlined,
  ExclamationCircleOutlined,
  CheckCircleOutlined,
  CloseCircleOutlined,
  PlusCircleOutlined,
  MinusCircleOutlined,
  LeftCircleOutlined,
  RightCircleOutlined,
  UpCircleOutlined,
  DownCircleOutlined,
  PlayCircleOutlined,
  PauseCircleOutlined,
  StopOutlined,
  FastForwardOutlined,
  FastBackwardOutlined,
  StepForwardOutlined,
  StepBackwardOutlined,
  CaretUpOutlined,
  CaretDownOutlined,
  CaretLeftOutlined,
  CaretRightOutlined,
  MenuOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  TableOutlined,
  LayoutOutlined,
  BlockOutlined,
  BorderOutlined,
  PictureOutlined,
  GiftOutlined,
  TrophyOutlined,
  CrownOutlined,
  FlagOutlined,
  CompassOutlined,
  WifiOutlined,
  MobileOutlined,
  TabletOutlined,
  LaptopOutlined,
  DesktopOutlined,
  CameraOutlined,
  VideoCameraOutlined,
  AudioOutlined,
  CustomerServiceOutlined,
  NotificationOutlined,
  SoundOutlined,
  BulbOutlined,
  ExperimentOutlined,
  HourglassOutlined,
  LoadingOutlined,
} from '@ant-design/icons-vue';

// 图标映射配置
export const iconMap: Record<string, any> = {
  // 基础图标
  dashboard: DashboardOutlined,
  home: HomeOutlined,
  user: UserOutlined,
  setting: SettingOutlined,
  team: TeamOutlined,
  cloud: CloudOutlined,
  file: FileOutlined,
  
  // 系统管理
  system: SettingOutlined,
  peoples: TeamOutlined,
  'tree-table': ApartmentOutlined,
  tree: ApartmentOutlined,
  dict: BookOutlined,
  edit: EditOutlined,
  
  // 监控模块
  monitor: MonitorOutlined,
  login: LoginOutlined,
  logininfor: LoginOutlined,
  form: FormOutlined,
  
  // 存储模块
  upload: UploadOutlined,
  database: DatabaseOutlined,
  server: CloudServerOutlined,
  documentation: FileTextOutlined,
  
  // 常用图标
  appstore: AppstoreOutlined,
  book: BookOutlined,
  bars: BarsOutlined,
  safety: SafetyOutlined,
  tool: ToolOutlined,
  folder: FolderOutlined,
  search: SearchOutlined,
  filter: FilterOutlined,
  reload: ReloadOutlined,
  sync: SyncOutlined,
  download: DownloadOutlined,
  export: ExportOutlined,
  import: ImportOutlined,
  printer: PrinterOutlined,
  copy: CopyOutlined,
  delete: DeleteOutlined,
  save: SaveOutlined,
  check: CheckOutlined,
  close: CloseOutlined,
  plus: PlusOutlined,
  minus: MinusOutlined,
  eye: EyeOutlined,
  lock: LockOutlined,
  unlock: UnlockOutlined,
  key: KeyOutlined,
  
  // 用户相关
  'user-add': UserAddOutlined,
  'usergroup-add': UsergroupAddOutlined,
  solution: SolutionOutlined,
  profile: ProfileOutlined,
  
  // 控制相关
  control: ControlOutlined,
  api: ApiOutlined,
  code: CodeOutlined,
  bug: BugOutlined,
  rocket: RocketOutlined,
  thunderbolt: ThunderboltOutlined,
  fire: FireOutlined,
  
  // 图表相关
  fund: FundOutlined,
  'bar-chart': BarChartOutlined,
  'line-chart': LineChartOutlined,
  'pie-chart': PieChartOutlined,
  
  // 商业相关
  'shopping-cart': ShoppingCartOutlined,
  shop: ShopOutlined,
  tag: TagOutlined,
  tags: TagsOutlined,
  
  // 通讯相关
  bell: BellOutlined,
  message: MessageOutlined,
  mail: MailOutlined,
  phone: PhoneOutlined,
  
  // 时间地点
  calendar: CalendarOutlined,
  'clock-circle': ClockCircleOutlined,
  environment: EnvironmentOutlined,
  global: GlobalOutlined,
  
  // 排序
  'sort-ascending': SortAscendingOutlined,
  'sort-descending': SortDescendingOutlined,
  
  // 状态图标
  star: StarOutlined,
  heart: HeartOutlined,
  like: LikeOutlined,
  dislike: DislikeOutlined,
  smile: SmileOutlined,
  frown: FrownOutlined,
  meh: MehOutlined,
  
  // 提示图标
  'question-circle': QuestionCircleOutlined,
  'info-circle': InfoCircleOutlined,
  warning: WarningOutlined,
  'exclamation-circle': ExclamationCircleOutlined,
  'check-circle': CheckCircleOutlined,
  'close-circle': CloseCircleOutlined,
  
  // 方向图标
  'left-circle': LeftCircleOutlined,
  'right-circle': RightCircleOutlined,
  'up-circle': UpCircleOutlined,
  'down-circle': DownCircleOutlined,
  'caret-up': CaretUpOutlined,
  'caret-down': CaretDownOutlined,
  'caret-left': CaretLeftOutlined,
  'caret-right': CaretRightOutlined,
  
  // 播放控制
  'play-circle': PlayCircleOutlined,
  'pause-circle': PauseCircleOutlined,
  stop: StopOutlined,
  'fast-forward': FastForwardOutlined,
  'fast-backward': FastBackwardOutlined,
  
  // 布局相关
  menu: MenuOutlined,
  'menu-fold': MenuFoldOutlined,
  'menu-unfold': MenuUnfoldOutlined,
  table: TableOutlined,
  layout: LayoutOutlined,
  block: BlockOutlined,
  border: BorderOutlined,
  
  // 媒体相关
  picture: PictureOutlined,
  camera: CameraOutlined,
  'video-camera': VideoCameraOutlined,
  audio: AudioOutlined,
  
  // 奖励相关
  gift: GiftOutlined,
  trophy: TrophyOutlined,
  crown: CrownOutlined,
  flag: FlagOutlined,
  
  // 设备相关
  mobile: MobileOutlined,
  tablet: TabletOutlined,
  laptop: LaptopOutlined,
  desktop: DesktopOutlined,
  wifi: WifiOutlined,
  
  // 其他
  compass: CompassOutlined,
  'customer-service': CustomerServiceOutlined,
  notification: NotificationOutlined,
  sound: SoundOutlined,
  bulb: BulbOutlined,
  experiment: ExperimentOutlined,
  hourglass: HourglassOutlined,
  loading: LoadingOutlined,
};

/**
 * 获取图标组件
 * @param iconName 图标名称
 * @returns 图标组件或null
 */
export const getIcon = (iconName?: string) => {
  if (!iconName) return null;
  const key = iconName.toLowerCase();
  return iconMap[key] || null;
};

/**
 * 获取所有可用图标列表（用于图标选择器）
 * @returns 图标名称数组
 */
export const getAvailableIcons = () => {
  return Object.keys(iconMap).sort();
};

/**
 * 检查图标是否存在
 * @param iconName 图标名称
 * @returns 是否存在
 */
export const hasIcon = (iconName?: string) => {
  if (!iconName) return false;
  return iconName.toLowerCase() in iconMap;
};
